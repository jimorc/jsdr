#include "jsdr_config.h"

#include <wx/gdicmn.h>

#include <cstdlib>
#include <fstream>
#include <locale>
#include <memory>
#include <string>

namespace jsdr {
   const char* const kConfigFileName = "jsdr.config";

   auto JSdrConfig::LoadDisplayProperties() -> ConfigFileStatus {
      const auto* home = getenv("HOME");   // NOLINT
      if (home != nullptr) {
         _configFileName = home + std::string(kConfigFileName);
      } else {
         return ConfigFileStatus::kNoUser;
      }
      if (ConfigFileExists()) {
         std::ifstream configStream;
         configStream.imbue(std::locale());
         configStream.open(_configFileName);
         configStream >> _values;
         configStream.close();
         return ConfigFileStatus::kOk;
      }
      SetDefaultConfigValues();
      return ConfigFileStatus::kFileInitialized;
   }

   auto JSdrConfig::StoreDisplayProperties() -> bool {
      if (_configFileName.empty()) { return false; }
      std::ofstream configStream;
      configStream.imbue(std::locale());
      configStream.open(_configFileName);
      configStream << _values << '\n';
      configStream.close();
      return true;
   }

   void JSdrConfig::SetDefaultDisplayValues() {
      _values["mainFrame"]["position"]["x"]  = wxDefaultPosition.x;
      _values["mainFrame"]["position"]["y"]  = wxDefaultPosition.y;
      _values["mainFrame"]["size"]["width"]  = wxDefaultSize.GetWidth();
      _values["mainFrame"]["size"]["height"] = wxDefaultSize.GetHeight();
   }

   void JSdrConfig::SetDefaultConfigValues() { SetDefaultDisplayValues(); }

   auto JSdrConfig::GetDisplayProperties() -> std::shared_ptr<JSdrConfig::DisplayProperties> {
      auto      displayProps          = std::make_shared<JSdrConfig::DisplayProperties>();
      const int x                     = _values["mainFrame"]["position"]["x"].asInt();   // NOLINT
      const int y                     = _values["mainFrame"]["position"]["y"].asInt();   // NOLINT
      displayProps->mainFramePosition = wxPoint(x, y);
      const int width                 = _values["mainFrame"]["size"]["width"].asInt();
      const int height                = _values["mainFrame"]["size"]["width"].asInt();
      displayProps->mainFrameSize     = wxSize(width, height);
      return displayProps;
   }
}   // namespace jsdr