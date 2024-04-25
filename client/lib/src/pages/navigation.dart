import 'package:client/src/pages/budget_page.dart';
import 'package:client/src/pages/me_page.dart';
import 'package:client/src/pages/home_page.dart';
import 'package:client/src/pages/transactions_page.dart';
import 'package:flutter/material.dart';

class Navigation extends StatefulWidget {
  const Navigation({super.key});

  @override
  State<Navigation> createState() => _NavigationState();
}

class _NavigationState extends State<Navigation> {
  int currentPageIndex = 0;

  @override
  Widget build(BuildContext context) {
    final ThemeData theme = Theme.of(context);
    return Scaffold(
      bottomNavigationBar: NavigationBar(
        onDestinationSelected: (int index) {
          setState(() {
            currentPageIndex = index;
          });
        },
        indicatorColor: Colors.blue,
        selectedIndex: currentPageIndex,
        destinations: const <NavigationDestination>[
          NavigationDestination(
            selectedIcon: Icon(Icons.home),
            icon: Icon(Icons.home_outlined),
            label: 'Home',
          ),
          NavigationDestination(
            selectedIcon: Icon(Icons.account_balance_wallet),
            icon: Icon(Icons.account_balance_wallet_outlined),
            label: 'Budget',
          ),
          NavigationDestination(
            selectedIcon: Icon(Icons.credit_card),
            icon: Icon(Icons.credit_card_outlined),
            label: 'Transactions',
          ),
          NavigationDestination(
            selectedIcon: Icon(Icons.person),
            icon: Icon(Icons.person_outline),
            label: 'Me',
          ),
        ],
      ),
      body: <Widget>[
        HomePage(), // Home page
        BudgetPage(), // Budget page
        TransactionsPage(), // Transactions page
        MePage(), // Help page
      ][currentPageIndex],
    );
  }
}