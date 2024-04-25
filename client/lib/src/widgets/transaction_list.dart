import 'package:client/src/models/transaction.dart';
import 'package:client/src/pages/auth_bloc.dart';
import 'package:client/src/utils/dio_service.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class TransactionList extends StatefulWidget {
  const TransactionList({super.key});

  @override
  _TransactionListState createState() => _TransactionListState();
}

class _TransactionListState extends State<TransactionList> {
  late final DioService dioService;
  late final AuthBloc authBloc;

  @override
  void initState() {
    super.initState();
    dioService = context.read<DioService>();
    authBloc = BlocProvider.of<AuthBloc>(context);
  }

  Future<List<Transaction>> fetchTransactions() async {
    String? accessToken;
    try {
      accessToken = await authBloc.getAccessToken();
    } catch (e) {
      print('Error getting access token: $e');
    }

    final response = await dioService.dio.get(
      'http://hisab-backend.eu-west-1.elasticbeanstalk.com/transactions',
      options: Options(
        headers: {
          'Authorization': 'Bearer $accessToken',
          'Content-Type': 'application/json',
        },
      ),
    );

    if (response.statusCode == 200) {
      final jsonData = response.data;

      return jsonData['transactions'].map<Transaction>((transactionData) {
        return Transaction(
          id: transactionData['id'],
          name: transactionData['name'],
          amount: double.parse(transactionData['amount']),
          dateTime: DateTime.parse(transactionData['date_time']),
          category: transactionData['category'],
          imageUrl: transactionData['image_url'],
        );
      }).toList();
    } else {
      throw Exception('Failed to fetch data');
    }
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<List<Transaction>>(
      future: fetchTransactions(),
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return Center(child: CircularProgressIndicator());
        } else if (snapshot.hasError) {
          return Center(child: Text('Error: ${snapshot.error}'));
        } else if (snapshot.hasData) {
          List<Transaction> transactions = snapshot.data!;
          return ListView.builder(
            itemCount: transactions.length,
            itemBuilder: (context, index) {
              final transaction = transactions[index];
              final isDateFirstTransaction = index == 0 ||
                  transaction.dateTime.day !=
                      transactions[index - 1].dateTime.day;
              Color textColor = transaction.amount >= 0 ? Colors.green : Colors.red;

              return Column(
                children: [
                  if (isDateFirstTransaction)
                    ListTile(
                      title: Text(
                        transaction.dateTime.toLocal().toString().split(' ')[0],
                        style: TextStyle(fontWeight: FontWeight.bold),
                      ),
                    ),
                  ListTile(
                    leading: _buildTransactionIcon(transaction.imageUrl),
                    title: Text(
                      transaction.name,
                      style: TextStyle(fontSize: 16),
                    ),
                    trailing: Text(
                      '${transaction.amount.toStringAsFixed(2)} AED',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                        color: textColor,
                      ),
                    ),
                  ),
                ],
              );
            },
          );
        } else {
          return Center(child: Text('No transactions found'));
        }
      },
    );
  }

  Widget _buildTransactionIcon(String imageUrl) {
    if (imageUrl.isNotEmpty) {
      return Image.network(
        imageUrl,
        width: 48,
        height: 48,
        errorBuilder: (context, error, stackTrace) {
          return Icon(Icons.shop);
        },
      );
    } else {
      return Icon(Icons.shop);
    }
  }
}
