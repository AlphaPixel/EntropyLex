#ifndef CPP_KOWALSKI_TESTS_MAIN_TEST_CC_
#define CPP_KOWALSKI_TESTS_MAIN_TEST_CC_

#include "util.h"
#include <gtest/gtest.h>

TEST(SchemaTest, FileOpensAsJson) {
  auto path = fixture_path("schema/el-8-valid.lxj");
  std::ifstream file(path);
  ASSERT_TRUE(file.is_open()) << "couldn't open " << path;

  nlohmann::json j;
  file >> j;
  EXPECT_EQ(j["profile"]["name"], "EL-8");
}

TEST(SchemaTest, FilePassesValidation) {
  auto schema = load_json_file(fixture_path("schema/lxj-v1.schema.json"));
  auto json = load_json_file(fixture_path("schema/el-8-valid.lxj"));

  nlohmann::json_schema::json_validator validator;
  ASSERT_NO_THROW(validator.set_root_schema(schema));

  EXPECT_NO_THROW(validator.validate(json));
}

TEST(SchemaTest, InvalidFileDoesNotPassValidation) {
  auto schema = load_json_file(fixture_path("schema/lxj-v1.schema.json"));
  auto json = load_json_file(fixture_path("schema/el-8-invalid.lxj"));

  nlohmann::json_schema::json_validator validator;

  EXPECT_THROW(validator.validate(json), std::exception);
}

#endif // CPP_KOWALSKI_TESTS_MAIN_TEST_CC_
