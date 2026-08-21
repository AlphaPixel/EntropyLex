#ifndef CPP_KOWALSKI_TESTS_MAIN_TEST_CC_
#define CPP_KOWALSKI_TESTS_MAIN_TEST_CC_

#include <gtest/gtest.h>

// Demonstrate some basic assertions.
TEST(HelloTest, BasicAssertions) {
  // Expect two strings not to be equal.
  EXPECT_STRNE("hello", "world");
  // Expect equality.
  EXPECT_EQ(7 * 6, 42);
}

#endif // CPP_KOWALSKI_TESTS_MAIN_TEST_CC_
