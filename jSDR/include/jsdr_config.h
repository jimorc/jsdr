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

   class jSDRConfig {
   public:
      struct DisplayProperties {
         wxPoint mainFramePosition = wxDefaultPosition;
         wxSize  mainFrameSize     = wxDefaultSize;
      };
      jSDRConfig()  = default;
      ~jSDRConfig() = default;
      ;
      auto         LoadDisplayProperties() -> ConfigFileStatus;
      auto         StoreDisplayProperties() -> bool;
      auto         GetDisplayProperties() -> std::shared_ptr<DisplayProperties>;
      Json::Value& values() { return m_values; }
      auto         ConfigFileExists() -> bool { return std::filesystem::exists(_configFileName); }
      void         SetDefaultConfigValues();
      void         SetDefaultDisplayValues();

   private:
      std::string _configFileName;
      Json::Value m_values;
   };
}   // namespace jsdr