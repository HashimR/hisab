import 'package:client/src/pages/auth_bloc.dart';
import 'package:lean_sdk_flutter/lean_sdk_flutter.dart';

class LeanService {

  final AuthBloc authBloc;

  LeanService(this.authBloc);

  Future<String?> getLeanCustomerId() async {
    String? leanCustomerId;
    try {
      leanCustomerId = await authBloc.getLeanCustomerId();
    } catch (e) {
      print('Error getting access token: $e');
    }
    return leanCustomerId;
  }

  List<LeanPermissions> getPermissions() {
    return [LeanPermissions.identity, LeanPermissions.accounts, LeanPermissions.balance, LeanPermissions.transactions];
  }

  String getAppToken() {
    return "INSERT_APP_TOKEN";
  }

}

