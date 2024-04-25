import 'dart:convert';
import 'package:client/src/pages/auth_bloc.dart';
import 'package:http/http.dart' as http;

class RegistrationService {
  static const String _baseUrl = 'http://hisab-backend.eu-west-1.elasticbeanstalk.com';
  final AuthBloc _authBloc;

  RegistrationService(this._authBloc);

  Future<void> registerUser({
    required String username,
    required String password,
    required String firstName,
    required String lastName,
    required String phoneNumber,
    required String country,
  }) async {
    final Uri registerUri = Uri.parse('$_baseUrl/register');

    final Map<String, String> requestBody = {
      'username': username,
      'password': password,
      'first_name': firstName,
      'last_name': lastName,
      'phone_number': phoneNumber,
      'country': country,
    };

    try {
      final response = await http.post(
        registerUri,
        body: jsonEncode(requestBody),
        headers: <String, String>{
          'Content-Type': 'application/json',
        },
      );

      if (response.statusCode == 201) {
        final Map<String, dynamic> responseData = jsonDecode(response.body);
        
        _authBloc.add(AuthUserRegistered(
          accessToken: responseData['access_token'],
          refreshToken: responseData['refresh_token'],
          leanCustomerId: responseData['lean_customer_id']
        ));
      } else if (response.statusCode == 400) {
        final Map<String, dynamic> errorData = jsonDecode(response.body);
        throw RegistrationException(errorData['message']);
      } else {
        throw RegistrationException('Registration failed. Please try again.');
      }
    } catch (e) {
      throw RegistrationException('Network error: $e');
    }
  }
}

class RegistrationException implements Exception {
  final String message;

  RegistrationException(this.message);
}
