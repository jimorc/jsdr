#ifndef JSDR_CONFIG_H
#define JSDR_CONFIG_H

#include <json/json.h>
#include <wx/wx.h>

#include <filesystem>
#include <memory>
#include <string>

namespace jsdr {
   enum class ConfigFileStatus {
      kOk,
      kNoUser,
      kFileInitialized,
   };

   class JSdrConfig {
   public:
      struct DisplayProperties {
         wxPoint mainFramePosition = wxDefaultPosition;
         wxSize  mainFrameSize     = wxDefaultSize;
      };
      JSdrConfig()                                   = default;
      JSdrConfig(const JSdrConfig&)                  = delete;
      JSdrConfig(JSdrConfig&&)                       = delete;
      auto operator=(const JSdrConfig) -> JSdrConfig = delete;
      auto operator=(JSdrConfig&&) -> JSdrConfig     = delete;
      ~JSdrConfig()                                  = default;
      ;
      auto LoadDisplayProperties() -> ConfigFileStatus;
      auto StoreDisplayProperties() -> bool;
      auto GetDisplayProperties() -> std::shared_ptr<DisplayProperties>;
      auto Values() noexcept -> Json::Value& { return _values; }
      auto ConfigFileExists() -> bool { return std::filesystem::exists(_configFileName); }
      void SetDefaultConfigValues();
      void SetDefaultDisplayValues();

   private:
      std::string _configFileName;
      Json::Value _values;
   };
}   // namespace jsdr

#endif   // JSDR_CONFIG_H