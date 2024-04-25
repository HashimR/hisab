import 'package:flutter/material.dart';
import 'package:intl_phone_field/countries.dart';
import 'package:intl_phone_field/intl_phone_field.dart';

class PhonePicker extends StatefulWidget {
  final TextEditingController phoneNumberController;

  const PhonePicker({
    super.key,
    required this.phoneNumberController,
  });

  @override
  _PhonePickerState createState() => _PhonePickerState();
}

class _PhonePickerState extends State<PhonePicker> {
  late Country _selectedCountry;
  String? _phoneNumberError;
  final String initialCountryCode = 'AE';

  @override
  void initState() {
    super.initState();
    _selectedCountry = countries
        .firstWhere((element) => element.code == initialCountryCode);
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        IntlPhoneField(
          autofocus: true,
          decoration: InputDecoration(
            labelText: 'Phone Number',
            filled: true,
            fillColor: Colors.white,
          ),
          initialCountryCode: initialCountryCode,
          onChanged: (phone) {
            final numberLength = phone.number.length;
            if (numberLength >= _selectedCountry.minLength &&
                numberLength <= _selectedCountry.maxLength) {
              setState(() {
                _phoneNumberError = null;
              });
              widget.phoneNumberController.text = phone.completeNumber;
            } else {
              setState(() {
                _phoneNumberError = 'Invalid phone number';
              });
            }
          },
          onCountryChanged: (country) {
            setState(() {
              _selectedCountry = country;
            });
          },
        ),
      ],
    );
  }
}
