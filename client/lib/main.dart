import 'package:client/src/pages/auth_bloc.dart';
import 'package:client/src/services/lean_service.dart';
import 'package:client/src/services/login_service.dart';
import 'package:client/src/services/refresh_service.dart';
import 'package:client/src/services/registration_service.dart';
import 'package:client/src/utils/dio_service.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:provider/provider.dart';

import 'src/app.dart';
import 'src/settings/settings_controller.dart';
import 'src/settings/settings_service.dart';

void main() async {
  // Set up the SettingsController, which will glue user settings to multiple
  // Flutter Widgets.
  final settingsController = SettingsController(SettingsService());

  await settingsController.loadSettings();
  WidgetsFlutterBinding.ensureInitialized();
  final authBloc = AuthBloc();
  final refreshService = RefreshService(authBloc);
  final loginService = LoginService(authBloc);
  final registrationService = RegistrationService(authBloc);
  final dioService =
      DioService(authBloc: authBloc, refreshService: refreshService);
  final leanService = LeanService(authBloc);

  runApp(MultiProvider(
    providers: [
      BlocProvider<AuthBloc>(create: (context) => authBloc),
      Provider<RefreshService>(create: (context) => refreshService),
      Provider<DioService>(create: (context) => dioService),
      Provider<LoginService>(create: (context) => loginService),
      Provider<RegistrationService>(create: (context) => registrationService),
      Provider<LeanService>(create: (context) => leanService),
    ],
    child: HisabApp(),
  ));
}
