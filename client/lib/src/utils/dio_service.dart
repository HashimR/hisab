import 'package:dio/dio.dart';
import 'package:client/src/pages/auth_bloc.dart';
import 'package:client/src/services/refresh_service.dart';
import 'package:jwt_decoder/jwt_decoder.dart';

class DioService {
  late Dio dio;
  final AuthBloc authBloc;
  final RefreshService refreshService;

  DioService({
    required this.authBloc,
    required this.refreshService,
  }) {
    dio = Dio();
    dio.interceptors.add(_tokenRefreshInterceptor());
  }

  InterceptorsWrapper _tokenRefreshInterceptor() {
    return InterceptorsWrapper(
      onRequest: (options, handler) async {
        await _handleAccessTokenExpiration(options);
        return handler.next(options);
      },
      onError: (DioException error, handler) async {
        if (error.response?.statusCode == 401) {
          await _retryWithRefreshToken(error, handler);
        } else {
          return handler.next(error);
        }
      },
    );
  }

  Future<void> _handleAccessTokenExpiration(RequestOptions options) async {
    final accessToken = await authBloc.getAccessToken();
    if (accessToken == null || JwtDecoder.isExpired(accessToken)) {
      await _attemptTokenRefresh();
    } else {
      options.headers['Authorization'] = 'Bearer $accessToken';
    }
  }

  Future<void> _attemptTokenRefresh() async {
    final refreshToken = await authBloc.getRefreshToken();
    if (refreshToken == null) {
      authBloc.add(AuthLogoutRequested());
      return;
      // throw Exception('No refresh token available');
    }

    try {
      await refreshService.refreshToken(refreshToken);
    } catch (e) {
      authBloc.add(AuthLogoutRequested());
      throw Exception('Token refresh failed');
    }
  }

  Future<void> _retryWithRefreshToken(DioException error, ErrorInterceptorHandler handler) async {
    try {
      await _attemptTokenRefresh();
      final newAccessToken = await authBloc.getAccessToken();
      if (newAccessToken != null) {
        error.requestOptions.headers['Authorization'] = 'Bearer $newAccessToken';
        final opts = Options(method: error.requestOptions.method, headers: error.requestOptions.headers);
        final cloneReq = await dio.request(
          error.requestOptions.path,
          options: opts,
          data: error.requestOptions.data,
          queryParameters: error.requestOptions.queryParameters,
        );
        handler.resolve(cloneReq);
      } else {
        handler.reject(error);
      }
    } catch (e) {
      authBloc.add(AuthLogoutRequested());
      handler.reject(error);
    }
  }
  Future<bool> isTokenExpired(String? accessToken) async {
    if (accessToken == null) {
      return true;
    }

    try {
      return JwtDecoder.isExpired(accessToken);
    } catch (e) {
      return true;
    }
  }
}
