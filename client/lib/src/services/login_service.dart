import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:client/src/pages/auth_bloc.dart';

class LoginService {
  static const String _baseUrl = 'http://localhost:8080';

  final AuthBloc authBloc;

  LoginService(this.authBloc);

  Future<void> loginUser({
    required String email,
    required String password,
  }) async {
    final Uri loginUri = Uri.parse('$_baseUrl/login');

    final Map<String, String> requestBody = {
      'email': email,
      'password': password,
    };

    try {
      final response = await http.post(
        loginUri,
        body: jsonEncode(requestBody),
        headers: <String, String>{
          'Content-Type': 'application/json',
        },
      );

      if (response.statusCode == 200) {
        final Map<String, dynamic> responseData = jsonDecode(response.body);
        final String accessToken = responseData['access_token'];
        final String refreshToken = responseData['refresh_token'];
        final String leanCustomerId = responseData['lean_customer_id'];
        final bool isConnected = responseData['is_connected'];
        authBloc.add(AuthLoginSuccessful(accessToken, refreshToken, leanCustomerId, isConnected));
      } else if (response.statusCode == 400) {
        final Map<String, dynamic> errorData = jsonDecode(response.body);
        // Handle specific error cases (e.g., validation errors)
        throw LoginException(errorData['message']);
      } else {
        // Handle other error cases with a generic error message
        throw LoginException('Login failed. Please try again.');
      }
    } catch (e) {
      // Handle network errors and other exceptions
      throw LoginException('Network error: $e');
    }
  }
}

class LoginException implements Exception {
  final String message;

  LoginException(this.message);
}
