import 'package:flutter/material.dart';
import '../landing-page.dart';

class OnboardingData {
  final String image;
  final String title;
  final String subtitle;

  const OnboardingData({
    required this.image,
    required this.title,
    required this.subtitle,
  });
}

class OnboardingScreen extends StatefulWidget {
  const OnboardingScreen({super.key});

  @override
  State<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends State<OnboardingScreen> {
  final PageController _pageController = PageController();
  int _currentPage = 0;

  static const List<OnboardingData> _onboardingPages = [
    OnboardingData(
      image: 'assets/images/onboarding_doctor.png',
      title: 'IMNCI Simplified',
      subtitle: 'A modern digital version of IMCI guidelines designed for daily clinical use.',
    ),
    OnboardingData(
      image: 'assets/images/onboarding_progress.png',
      title: 'Guided Clinical Assessment',
      subtitle: 'Instantly identify mild, moderate, or severe conditions using IMCI logic.',
    ),
    OnboardingData(
      image: 'assets/images/onboarding_team.png',
      title: 'Treatment & Monitoring',
      subtitle: 'Support proper treatment and ensure timely follow-up care.',
    ),
  ];

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  void _onNext() {
    if (_currentPage < _onboardingPages.length - 1) {
      _pageController.nextPage(
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeInOut,
      );
    } else {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (context) => const LandingPage()),
      );
    }
  }

  void _onBack() {
    if (_currentPage > 0) {
      _pageController.previousPage(
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeInOut,
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final size = MediaQuery.of(context).size;

    return Scaffold(
      backgroundColor: Colors.white,
      body: SafeArea(
        child: Column(
          children: [
            Expanded(
              child: PageView.builder(
                controller: _pageController,
                itemCount: _onboardingPages.length,
                onPageChanged: (index) {
                  setState(() {
                    _currentPage = index;
                  });
                },
                itemBuilder: (context, index) {
                  final data = _onboardingPages[index];
                  return Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 24.0, vertical: 40.0),
                    child: Column(
                      children: [
                        const Spacer(flex: 2),
                        
                        // Image
                        Center(
                          child: Image.asset(
                            data.image,
                            height: size.height * 0.4,
                            fit: BoxFit.contain,
                          ),
                        ),
                        
                        const Spacer(flex: 1),
                        
                        // Title
                        Text(
                          data.title,
                          textAlign: TextAlign.center,
                          style: theme.textTheme.headlineSmall?.copyWith(
                            fontWeight: FontWeight.bold,
                            color: const Color(0xFF2D4373),
                            fontSize: 28,
                          ),
                        ),
                        
                        const SizedBox(height: 16),
                        
                        // Subtitle
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 20.0),
                          child: Text(
                            data.subtitle,
                            textAlign: TextAlign.center,
                            style: theme.textTheme.bodyLarge?.copyWith(
                              color: Colors.blueGrey[400],
                              height: 1.5,
                            ),
                          ),
                        ),
                        
                        const Spacer(flex: 3),
                      ],
                    ),
                  );
                },
              ),
            ),
            
            // Bottom Section: Indicators and Buttons
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24.0, vertical: 32.0),
              child: Stack(
                alignment: Alignment.center,
                children: [
                  // Page Indicator
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: List.generate(
                      _onboardingPages.length,
                      (index) => AnimatedContainer(
                        duration: const Duration(milliseconds: 300),
                        margin: const EdgeInsets.only(right: 8),
                        width: _currentPage == index ? 24 : 8,
                        height: 8,
                        decoration: BoxDecoration(
                          color: _currentPage == index 
                              ? const Color(0xFF5890D8) 
                              : const Color(0xFFE0E0E0),
                          borderRadius: BorderRadius.circular(4),
                        ),
                      ),
                    ),
                  ),

                  // Buttons
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      // Back Button (only if not on first page)
                      if (_currentPage > 0)
                        OutlinedButton(
                          onPressed: _onBack,
                          style: OutlinedButton.styleFrom(
                            side: const BorderSide(color: Color(0xFF5890D8), width: 1.5),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(16),
                            ),
                            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                          ),
                          child: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              const Icon(
                                Icons.chevron_left,
                                color: Color(0xFF5890D8),
                                size: 20,
                              ),
                              const SizedBox(width: 8),
                              const Text(
                                'Back',
                                style: TextStyle(
                                  color: Color(0xFF5890D8),
                                  fontWeight: FontWeight.bold,
                                  fontSize: 16,
                                ),
                              ),
                            ],
                          ),
                        )
                      else
                        const SizedBox(width: 60), // Placeholder to keep Next button right-aligned

                      // Next / Get Started Button
                      OutlinedButton(
                        onPressed: _onNext,
                        style: OutlinedButton.styleFrom(
                          side: const BorderSide(color: Color(0xFF5890D8), width: 1.5),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(16),
                          ),
                          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(
                              _currentPage == _onboardingPages.length - 1 
                                  ? 'Get Started' 
                                  : 'Next',
                              style: const TextStyle(
                                color: Color(0xFF5890D8),
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                              ),
                            ),
                            if (_currentPage < _onboardingPages.length - 1) ...[
                              const SizedBox(width: 8),
                              const Icon(
                                Icons.chevron_right,
                                color: Color(0xFF5890D8),
                                size: 20,
                              ),
                            ],
                          ],
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
