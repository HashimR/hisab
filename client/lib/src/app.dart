import 'package:client/src/pages/auth_bloc.dart';
import 'package:client/src/pages/initial_connection.dart';
import 'package:client/src/pages/navigation.dart';
import 'package:client/src/pages/welcome_page.dart';
import 'package:client/src/services/lean_service.dart';
import 'package:client/src/services/refresh_service.dart';
import 'package:client/src/utils/dio_service.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class HisabApp extends StatefulWidget {
  const HisabApp({super.key});

  @override
  State<HisabApp> createState() => _HisabAppState();
}

class _HisabAppState extends State<HisabApp> {
  late final AuthBloc authBloc;
  late final DioService dioService;
  late final RefreshService refreshService;
  late final LeanService leanService;

  @override
  void initState() {
    super.initState();
    authBloc = context.read<AuthBloc>()..add(AuthCheckRequested());
    leanService = context.read<LeanService>();
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(useMaterial3: true),
      home: StreamBuilder<AuthState>(
        stream: authBloc.stream,
        builder: (context, snapshot) {
          if (!snapshot.hasData) {
            return CircularProgressIndicator(); // Show loading indicator while waiting for initial state
          }

          final state = snapshot.data;

          // TODO: Add UnconfirmedConnectionState
          if (state is UnauthenticatedState) {
            return AuthenticationPage(); // Show authentication page by default
          } else if (state is UnconnectedState) {
            return InitialConnectionPage(leanService: leanService);
          } else if (state is AuthenticatedState) {
            return Navigation(); // Navigate to the authenticated page
          }
          // Return placeholder
          return Container();
        },
      ),
    );
  }
}
