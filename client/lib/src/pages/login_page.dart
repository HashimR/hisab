import 'package:client/src/pages/auth_bloc.dart';
import 'package:client/src/pages/initial_connection.dart';
import 'package:client/src/pages/navigation.dart';
import 'package:client/src/pages/welcome_page.dart';
import 'package:client/src/services/lean_service.dart';
import 'package:client/src/services/login_service.dart';
import 'package:email_validator/email_validator.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  _LoginPageState createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
  late final LoginService _loginService;
  late final LeanService _leanService;

  late TextEditingController _usernameController;
  late TextEditingController _passwordController;

  @override
  void initState() {
    super.initState();
    _usernameController = TextEditingController();
    _passwordController = TextEditingController();
    _loginService = context.read<LoginService>();
    _leanService = context.read<LeanService>();
  }

  @override
  void dispose() {
    // Dispose of the controllers to prevent memory leaks
    _usernameController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Color(0xFF49BDD1),
      appBar: AppBar(
        title: Text('Log In'),
      ),
      body: SingleChildScrollView(
        child: Form(
          key: _formKey,
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              children: [
                // Username TextField
                TextFormField(
                  enableSuggestions: false,
                  controller: _usernameController,
                  decoration: InputDecoration(
                    labelText: 'Email Address',
                    filled: true,
                    fillColor: Colors.white,
                  ),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Please enter your email address';
                    }
                    if (!EmailValidator.validate(value)) {
                      return 'Please enter a valid email address';
                    }
                    return null;
                  },
                ),
                SizedBox(height: MediaQuery.of(context).size.height * 0.02),

                // Password TextField
                TextFormField(
                  controller: _passwordController,
                  obscureText: true,
                  decoration: InputDecoration(
                    labelText: 'Password',
                    filled: true,
                    fillColor: Colors.white,
                  ),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Please enter your password';
                    }
                    return null;
                  },
                ),

                SizedBox(height: MediaQuery.of(context).size.height * 0.02),

                ElevatedButton(
                  onPressed: () async {
                    if (_formKey.currentState!.validate()) {
                      _formKey.currentState!.save();
                      try {
                        // final loginService = LoginService(_authBloc);

                        await _loginService.loginUser(
                          email: _usernameController.text,
                          password: _passwordController.text,
                        );
                        StatefulWidget nextPage = getNextPage(_loginService.authBloc);
                        Navigator.pushReplacement(
                          context,
                          MaterialPageRoute(builder: (context) => nextPage),
                        );
                      } catch (e) {
                        // Handle registration error and show the error message in a popup
                        String errorMessage = 'An error occurred during login';
                        if (e is LoginException) {
                          errorMessage = e.message;
                        }
                        showErrorPopup(context, errorMessage);
                      }
                    }
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.white,
                    foregroundColor: Colors.black,
                  ),
                  child: Text('Log In'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
  
  getNextPage(AuthBloc authBloc) {
    AuthState state = authBloc.state;
    if (authBloc.state is UnauthenticatedState || state is UnauthenticatedState) {
      return AuthenticationPage(); // Show welcome page by default
    } else if (state is UnconnectedState) {
      return InitialConnectionPage(leanService: _leanService);
    } else if (state is AuthenticatedState) {
      return Navigation(); // Navigate to the authenticated page
    }
    return Container();
  }
}

void showErrorPopup(BuildContext context, String errorMessage) {
  showDialog(
    context: context,
    builder: (BuildContext context) {
      return AlertDialog(
        title: Text('Error'),
        content: Text(errorMessage),
        actions: <Widget>[
          TextButton(
            onPressed: () {
              Navigator.of(context).pop();
            },
            child: Text('OK'),
          ),
        ],
      );
    },
  );
}
