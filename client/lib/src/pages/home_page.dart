import 'package:client/src/widgets/money_chart.dart';
import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';

class HomePage extends StatelessWidget {
  const HomePage({Key? key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Home'),
      ),
      body: Column(
        children: [
          // Line Chart
          Expanded(
            flex: 1,
            child: Container(
              color: Colors.grey[200], // Light blue background color
              child: MoneyChart()
            ),
          ),
           Expanded(
            flex: 1, // Use the same flex value to maintain a 50/50 split
            child: Container(
              color: Colors.grey[200], // Placeholder background color
              child: Center(
                child: Text(
                  'Placeholder for Bottom Half',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}