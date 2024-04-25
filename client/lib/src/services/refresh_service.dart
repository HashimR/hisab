import 'dart:convert';
import 'package:client/src/pages/auth_bloc.dart';
import 'package:http/http.dart' as http;

class RefreshService {
  static const String _baseUrl = 'http://hisab-backend.eu-west-1.elasticbeanstalk.com'; 

  final AuthBloc authBloc;

  RefreshService(this.authBloc);

  Future<void> refreshToken(String refreshToken) async {
    final Uri refreshTokenUri = Uri.parse('$_baseUrl/refresh');

    final Map<String, String> requestBody = {
      'refresh_token': refreshToken,
    };

    try {
      final response = await http.post(
        refreshTokenUri,
        body: jsonEncode(requestBody),
        headers: <String, String>{
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $refreshToken',
        },
      );

      if (response.statusCode == 200) {
        final Map<String, dynamic> responseData = jsonDecode(response.body);
        final String newAccessToken = responseData['access_token'];
        final String newRefreshToken = responseData['refresh_token'];
        authBloc.add(AuthRefreshSuccessful(newAccessToken, newRefreshToken));
      } else {
        // Handle error cases here and throw an exception
        throw RefreshException('Failed to refresh token');
      }
    } catch (e) {
      // Handle network errors and other exceptions and throw an exception
      throw RefreshException('Network error: $e');
    }
  }
}

class RefreshException implements Exception {
  final String message;

  RefreshException(this.message);
}