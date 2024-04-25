import 'package:client/src/pages/auth_bloc.dart';
import 'package:client/src/pages/navigation.dart';
import 'package:client/src/pages/welcome_page.dart';
import 'package:client/src/services/lean_service.dart';
import 'package:flutter/material.dart';
import 'package:lean_sdk_flutter/lean_sdk_flutter.dart';

class InitialConnectionPage extends StatefulWidget {
  final LeanService leanService;

  const InitialConnectionPage({super.key, required this.leanService});

  @override
  _InitialConnectionPageState createState() => _InitialConnectionPageState();
}

class _InitialConnectionPageState extends State<InitialConnectionPage> {
  late final LeanService _leanService;
  late final AuthBloc _authBloc;
  late String customerId;

  @override
  void initState() {
    super.initState();
    _leanService = widget.leanService;
    _authBloc = _leanService.authBloc;
    _leanService.getLeanCustomerId().then((String? id) {
      setState(() {
        customerId = id ?? '';
      });
    }).catchError((error) {});
  }

  @override
  Widget build(BuildContext context) {
    connect() {
      showModalBottomSheet(
          isScrollControlled: true,
          backgroundColor: Colors.transparent,
          context: context,
          builder: (context) {
            return Padding(
                padding: EdgeInsets.only(
                    bottom: MediaQuery.of(context).viewInsets.bottom),
                child: SizedBox(
                    height: MediaQuery.of(context).size.height * 0.8,
                    child: Lean.connect(
                      showLogs: true,
                      isSandbox: true,
                      country: LeanCountry.ae,
                      language: LeanLanguage.en,
                      appToken: _leanService.getAppToken(),
                      customerId: customerId,
                      permissions: _leanService.getPermissions(),
                      customization: const {
                        "button_text_color": "white",
                        "theme_color": "red",
                        "button_border_radius": "10",
                        "overlay_color": "pink",
                      },
                      callback: (resp) {
                        print("Callback: $resp");
                        if (resp.status == "SUCCESS") {
                          Navigator.pushReplacement(
                            context,
                            MaterialPageRoute(
                                builder: (context) => Navigation()),
                          );
                        }
                      },
                      actionCancelled: () => Navigator.pop(context),
                    )));
          });
    }

    return Scaffold(
      backgroundColor: Color(0xFF49BDD1),
      appBar: AppBar(
        title: Text('Connect Your Bank'),
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            children: [
              ElevatedButton(
                  onPressed: () {
                    _authBloc.add(AuthLogoutRequested());
                    // Navigate to the WelcomePage after logging out
                    Navigator.pushAndRemoveUntil(
                      context,
                      MaterialPageRoute(builder: (context) => AuthenticationPage()),
                      (route) => false, // Remove all previous routes
                    );
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.white,
                    foregroundColor: Colors.black,
                    padding: EdgeInsets.symmetric(
                        horizontal: 40, vertical: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius:
                          BorderRadius.circular(30.0),
                    ),
                  ),
                  child: Text(
                    'Log Out',
                    style: TextStyle(fontSize: 18),
                  ),
              ),
              ElevatedButton(
                onPressed: () => connect(),
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.white,
                  foregroundColor: Colors.black,
                ),
                child: Text('Connect'),
              ),
            ],
          ),
        ),
      ),
    );
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
