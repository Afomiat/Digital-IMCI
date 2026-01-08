import 'package:flutter/material.dart';
import 'dart:async';
import 'package:flutter_svg/flutter_svg.dart';

class LandingPage extends StatefulWidget {
  const LandingPage({super.key});

  @override
  State<LandingPage> createState() => _LandingPageState();
}

class _LandingPageState extends State<LandingPage>
    with SingleTickerProviderStateMixin {
  late AnimationController _animationController;
  late Animation<double> _fadeAnimation;

  @override
  void initState() {
    super.initState();
    
    // Initialize animation controller
    _animationController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1500),
    );

    // Create fade animation
    _fadeAnimation = Tween<double>(
      begin: 0.0,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.easeIn,
    ));

    // Start animation
    _animationController.forward();

    // Navigate to next screen after delay
    _navigateToNext();
  }

  void _navigateToNext() {
    Timer(const Duration(seconds: 3), () {
      if (mounted) {
        // TODO: Replace with your actual home/landing screen
        // Navigator.pushReplacement(
        //   context,
        //   MaterialPageRoute(builder: (context) => const HomeScreen()),
        // );
      }
    });
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      backgroundColor: Colors.white,
      body: SafeArea(
        child: FadeTransition(
          opacity: _fadeAnimation,
          child: Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                // App Logo/Icon
                Container(
                  // width: 120,
                  // height: 120,
                  // padding: const EdgeInsets.all(20), // Adds padding so the logo implies the circle shape nicely
                  // decoration: BoxDecoration(
                  //   color: colorScheme.onPrimary,
                  //   shape: BoxShape.circle,
                  //   boxShadow: [
                  //     BoxShadow(
                  //       color: Colors.black.withOpacity(0.2),
                  //       blurRadius: 20,
                  //       offset: const Offset(0, 10),
                  //     ),
                  //   ],
                  // ),
                  child: SvgPicture.asset(
                    'assets/images/white-logo.svg',
                    // Since the background (onPrimary) is likely white, 
                    // you might need to tint the white logo to be visible.
                    // If you want it to match the old icon color (primary), keep this line:
                    colorFilter: ColorFilter.mode(Color(0xFF5890D8), BlendMode.srcIn),
                  ),
                ),
                const SizedBox(height: 40),

                // App Name
                // Text(
                //   'IMNCI',
                //   style: theme.textTheme.displayLarge?.copyWith(
                //     color: Color(0xFF5890D8),
                //     fontWeight: FontWeight.bold,
                //     letterSpacing: 2,
                //   ),
                // ),
                // const SizedBox(height: 12),

                // Subtitle
                Text(
                  'Integrated Management of\nNewborn and Childhood Illness',
                  textAlign: TextAlign.center,
                  style: theme.textTheme.titleLarge?.copyWith(
                    color: Color(0xFF5890D8),
                    letterSpacing: 0.5,
                    height: 1.5,
                  ),
                ),
                const SizedBox(height: 60),

                // Loading Indicator
                // SizedBox(
                //   width: 40,
                //   height: 40,
                //   child: CircularProgressIndicator(
                //     valueColor: AlwaysStoppedAnimation<Color>(
                //       colorScheme.onPrimary,
                //     ),
                //     strokeWidth: 3,
                //   ),
                // ),
                // SizedBox(height: 12),
                Container(
                //   padding: const EdgeInsets.only(left: 50.0),
                  child: Image.asset(
                    // 'assets/images/couple-walks-with-baby.svg',
                    'assets/images/child-with-doctor.png',
                    // colorFilter: ColorFilter.mode(Colors.white, BlendMode.srcIn),
                    height: 200,
                    width: 200,
                  ),
                ),
                const SizedBox(height: 40),
                ElevatedButton(
                  onPressed: () {
                    // TODO: Navigate to signup page for invited health professionals
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF5890D8),
                    foregroundColor: Colors.white,
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  child: Text(
                    'Get Started',
                    style: theme.textTheme.labelLarge?.copyWith(
                      color: Colors.white,
                      letterSpacing: 0.5,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

