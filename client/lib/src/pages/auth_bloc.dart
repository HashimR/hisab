import 'dart:async';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final FlutterSecureStorage secureStorage = const FlutterSecureStorage(
    iOptions: IOSOptions(accessibility: KeychainAccessibility.first_unlock),
  );
  static const String accessTokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String leanCustomerIdKey = 'lean_customer_id';
  static const String isConnectedOnDeviceKey = 'is_connected_on_device';
  static const String isConfirmedConnectedKey = 'is_confirmed_connected';

  AuthBloc() : super(UnauthenticatedState()) {
    on<AuthUserRegistered>(_onUserRegistered);
    on<AuthCheckRequested>((event, emit) async {
      AuthState currentState = await _getCurrentState();
      emit(currentState);
    });
    on<AuthLogoutRequested>((event, emit) async {
      await deleteTokens();
      emit(UnauthenticatedState());
    });
    on<AuthLoginSuccessful>(_onUserLogin);
    on<AuthConnectedOnDevice>(_onConnectedOnDevice);
    on<AuthConfirmedConnected>(_onConfirmedConnected);
    on<AuthRefreshSuccessful>(_onAuthRefresh);
  }

  Future<AuthState> _getCurrentState() async {
    final accessToken = await getAccessToken();
    final refreshToken = await getRefreshToken();
    final leanCustomerId = await getLeanCustomerId();
    final isConnectedOnDevice = await getIsConnectedOnDevice();
    final isConfirmedConnected = await getIsConfirmedConnected();
    if (accessToken == null || refreshToken == null || leanCustomerId == null) {
      return UnauthenticatedState();
    } else if (!isConnectedOnDevice && !isConfirmedConnected) {
      return UnconnectedState(accessToken, refreshToken, leanCustomerId);
    } else if (!isConfirmedConnected) {
      return UnconfirmedConnectionState(accessToken, refreshToken, leanCustomerId);
    } else {
      return AuthenticatedState(accessToken, refreshToken, leanCustomerId);
    }
  }

  Future<void> _onUserLogin(
      AuthLoginSuccessful event, Emitter<AuthState> emit) async {
    await storeTokens(
      accessToken: event.accessToken,
      refreshToken: event.refreshToken,
      leanCustomerId: event.leanCustomerId
    );
    if (event.isConfirmedConnected) {
      await storeIsConfirmedConnected();
    }
    AuthState currentState = await _getCurrentState();
    emit(currentState);
  }

  Future<void> _onAuthRefresh(
      AuthRefreshSuccessful event, Emitter<AuthState> emit) async {
    String? leanCustomerId = await getLeanCustomerId();
    await storeTokens(
      accessToken: event.accessToken,
      refreshToken: event.refreshToken,
      leanCustomerId: leanCustomerId!
    );
    AuthState currentState = await _getCurrentState();
    emit(currentState);
  }

  Future<void> _onUserRegistered(
      AuthUserRegistered event, Emitter<AuthState> emit) async {
    await storeTokens(
      accessToken: event.accessToken,
      refreshToken: event.refreshToken,
      leanCustomerId: event.leanCustomerId
    );
    AuthState currentState = await _getCurrentState();
    emit(currentState);
  }

  Future<void> _onConnectedOnDevice(
      AuthConnectedOnDevice event, Emitter<AuthState> emit) async {
    storeIsConnectedOnDevice();
    AuthState currentState = await _getCurrentState();
    emit(currentState);
  }

  Future<void> _onConfirmedConnected(
      AuthConfirmedConnected event, Emitter<AuthState> emit) async {  
    storeIsConnectedOnDevice();
    storeIsConfirmedConnected();
    AuthState currentState = await _getCurrentState();
    emit(currentState);
  }

  // Storage
  Future<void> storeTokens({
    required String accessToken,
    required String refreshToken,
    required String leanCustomerId
  }) async {
    await secureStorage.write(key: accessTokenKey, value: accessToken);
    await secureStorage.write(key: refreshTokenKey, value: refreshToken);
    await secureStorage.write(key: leanCustomerIdKey, value: leanCustomerId);
  }

  Future<void> storeIsConnectedOnDevice() async {
    await secureStorage.write(key: isConnectedOnDeviceKey, value: "true");
  }

  Future<void> storeIsConfirmedConnected() async {
    await secureStorage.write(key: isConfirmedConnectedKey, value: "true");
  }

  // Getters
  Future<String?> getAccessToken() async {
    return await secureStorage.read(key: accessTokenKey);
  }

  Future<String?> getRefreshToken() async {
    return await secureStorage.read(key: refreshTokenKey);
  }

  Future<String?> getLeanCustomerId() async {
    return secureStorage.read(key: leanCustomerIdKey);
  }

  Future<bool> getIsConnectedOnDevice() async {
    return secureStorage.containsKey(key: isConnectedOnDeviceKey);
  }

  Future<bool> getIsConfirmedConnected() async {
    return secureStorage.containsKey(key: isConfirmedConnectedKey);
  }

  Future<void> deleteTokens() async {
    await secureStorage.delete(key: accessTokenKey);
    await secureStorage.delete(key: refreshTokenKey);
    await secureStorage.delete(key: leanCustomerIdKey);
    await secureStorage.delete(key: isConnectedOnDeviceKey);
    await secureStorage.delete(key: isConfirmedConnectedKey);
  }
}

// Auth events
abstract class AuthEvent {}

class AuthUserRegistered extends AuthEvent {
  final String accessToken;
  final String refreshToken;
  final String leanCustomerId;

  AuthUserRegistered({required this.accessToken, required this.refreshToken, required this.leanCustomerId});
}

class AuthUserRegisteredError extends AuthEvent {
  final String errorMessage;

  AuthUserRegisteredError(this.errorMessage);
}

class AuthCheckRequested extends AuthEvent {}

class AuthConnectedOnDevice extends AuthEvent {}

class AuthConfirmedConnected extends AuthEvent {}

class AuthLogoutRequested extends AuthEvent {}

class AuthLoginSuccessful extends AuthEvent {
  final String accessToken;
  final String refreshToken;
  final String leanCustomerId;
  final bool isConfirmedConnected;

  AuthLoginSuccessful(this.accessToken, this.refreshToken, this.leanCustomerId, this.isConfirmedConnected);
}

class AuthRefreshSuccessful extends AuthEvent {
  final String accessToken;
  final String refreshToken;

  AuthRefreshSuccessful(this.accessToken, this.refreshToken);
}

class AuthRefreshTokenRequested extends AuthEvent {}

// Auth states
abstract class AuthState {}

class UnauthenticatedState extends AuthState {}

// Authenticated from Hisab backend, but no entities connected
class UnconnectedState extends AuthState {
  String accessToken;
  String refreshToken;
  String leanCustomerId;

  UnconnectedState(this.accessToken, this.refreshToken, this.leanCustomerId);
}

// User has connected on client, but backend has not received data from Lean yet
class UnconfirmedConnectionState extends AuthState {
  String accessToken;
  String refreshToken;
  String leanCustomerId;

  UnconfirmedConnectionState(this.accessToken, this.refreshToken, this.leanCustomerId);
}

// User has connected and backend has Lean data
class AuthenticatedState extends AuthState {
  String accessToken;
  String refreshToken;
  String leanCustomerId; 

  AuthenticatedState(this.accessToken, this.refreshToken, this.leanCustomerId);
}

// class AuthRegistrationErrorState extends AuthState {
//   final String errorMessage;

//   AuthRegistrationErrorState(this.errorMessage);
// }




