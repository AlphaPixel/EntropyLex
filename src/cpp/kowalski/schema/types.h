#ifndef TYPES_H_
#define TYPES_H_
#include <cstdint>
#include <nlohmann/detail/macro_scope.hpp>
#include <nlohmann/json.hpp>
#include <string>
#include <vector>

struct ELX_Format {
  std::string name;
  uint32_t version;
  NLOHMANN_DEFINE_TYPE_INTRUSIVE(ELX_Format, name, version)
};

struct ELX_Dictionary {
  std::string id;
  uint32_t version;
  std::string purpose;
  std::string name;
  std::string language;
  std::string script;
  NLOHMANN_DEFINE_TYPE_INTRUSIVE(ELX_Dictionary, id, version, purpose, name,
                                 language, script)
};

struct ELX_Profile {
  std::string name;
  uint32_t normal_bits;
  std::vector<uint32_t> remainder_bits;
  uint32_t normal_count;
  uint32_t trim_count;
  uint32_t total_count;
  NLOHMANN_DEFINE_TYPE_INTRUSIVE(ELX_Profile, name, normal_bits, remainder_bits,
                                 normal_count, trim_count, total_count)
};

struct ELX_Tokenization {
  std::string kind;
  std::string canonical_separator;
  std::vector<std::string> separators;
  NLOHMANN_DEFINE_TYPE_INTRUSIVE(ELX_Tokenization, kind, canonical_separator,
                                 separators)
};

struct ELX_Recognition {
  std::string unicode_version;
  std::string normalization;
  std::string r_case;
  std::string token_text;
  ELX_Tokenization tokenization;
  NLOHMANN_DEFINE_TYPE_INTRUSIVE(ELX_Recognition, unicode_version,
                                 normalization, r_case, token_text,
                                 tokenization)
};

struct ELX_Fingerprint {
  std::string recipe;
  std::string sha256;
  NLOHMANN_DEFINE_TYPE_INTRUSIVE(ELX_Fingerprint, recipe, sha256)
};

struct ELX_Source {
  std::string name;
  std::string version;
  std::string role;
  std::string location;
  std::string license;
  std::string sha256;
  uint32_t record_count;
  NLOHMANN_DEFINE_TYPE_INTRUSIVE(ELX_Source, name, version, role, location,
                                 license, sha256, record_count)
};
struct ELX_Selection {
  std::string method;
  std::string settings_location;
  std::string settings_sha256;
  NLOHMANN_DEFINE_TYPE_INTRUSIVE(ELX_Selection, method, settings_location,
                                 settings_sha256)
};

struct ELX_Provenance {
  std::vector<ELX_Source> sources;
  ELX_Selection selection;
  NLOHMANN_DEFINE_TYPE_INTRUSIVE(ELX_Provenance, sources, selection)
};

class ELXJ {

public:
  ELX_Format format;
  ELX_Dictionary dictionary;
  ELX_Profile profile;
  ELX_Recognition recognition;
  std::string indexing;
  std::vector<std::string> tokens;
  NLOHMANN_DEFINE_TYPE_INTRUSIVE(ELXJ, format, dictionary, recognition,
                                 indexing, tokens)
};

#endif // TYPES_H_
