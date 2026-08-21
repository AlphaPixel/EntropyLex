#ifndef UTIL_H_
#define UTIL_H_
#include <filesystem>
#include <fstream>
#include <nlohmann/json-schema.hpp>
#include <nlohmann/json.hpp>

std::filesystem::path fixture_path(const std::string &rel) {
  return std::filesystem::path(PROJECT_SOURCE_DIR) / rel;
}
nlohmann::json load_json_file(const std::string &path) {
  std::ifstream file(path);
  if (!file.is_open()) {
    throw std::runtime_error("could not open file: " + path);
  }
  nlohmann::json j;
  file >> j; // parses directly from the stream
  return j;
}

#endif // UTIL_H_
