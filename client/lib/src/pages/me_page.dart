import 'package:client/src/pages/auth_bloc.dart';
import 'package:client/src/pages/welcome_page.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class MePage extends StatelessWidget {
  const MePage({super.key});

  @override
  Widget build(BuildContext context) {
    final AuthBloc authBloc = BlocProvider.of<AuthBloc>(context);

    return Scaffold(
      backgroundColor: Color(0x00000000),
      body: Column(
        mainAxisAlignment:
            MainAxisAlignment.center, // Center the children vertically
        children: [
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              children: [
                ElevatedButton(
                  onPressed: () {
                    authBloc.add(AuthLogoutRequested());
                    // Navigate to the WelcomePage after logging out
                    Navigator.pushAndRemoveUntil(
                      context,
                      MaterialPageRoute(builder: (context) => AuthenticationPage()),
                      (route) => false, // Remove all previous routes
                    );
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.white, // Button background color
                    foregroundColor: Colors.black, // Text color
                    padding: EdgeInsets.symmetric(
                        horizontal: 40, vertical: 16), // Button padding
                    shape: RoundedRectangleBorder(
                      borderRadius:
                          BorderRadius.circular(30.0), // Button border radius
                    ),
                  ),
                  child: Text(
                    'Log Out',
                    style: TextStyle(fontSize: 18),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
