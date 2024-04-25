class Transaction {
  final int id;
  final String name;
  final double amount;
  final DateTime dateTime;
  final String category;
  final String imageUrl;

  Transaction({
    required this.id,
    required this.name,
    required this.amount,
    required this.dateTime,
    required this.category,
    required this.imageUrl,
  });
}