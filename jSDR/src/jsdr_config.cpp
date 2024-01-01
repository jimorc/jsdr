#include "jsdr_config.h"

#include <wx/gdicmn.h>

#include <cstdlib>
#include <fstream>
#include <locale>
#include <memory>
#include <string>

namespace jsdr {
   const char* const kConfigFileName = "jsdr.config";

   auto jSDRConfig::LoadDisplayProperties() -> ConfigFileStatus {
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
         configStream >> m_values;
         configStream.close();
         return ConfigFileStatus::kOk;
      }
      SetDefaultConfigValues();
      return ConfigFileStatus::kFileInitialized;
   }

   auto jSDRConfig::StoreDisplayProperties() -> bool {
      if (_configFileName.empty()) { return false; }
      std::ofstream configStream;
      configStream.imbue(std::locale());
      configStream.open(_configFileName);
      configStream << m_values << '\n';
      configStream.close();
      return true;
   }

   void jSDRConfig::SetDefaultDisplayValues() {
      m_values["mainFrame"]["displayNumber"]  = 0;
      m_values["mainFrame"]["position"]["x"]  = wxDefaultPosition.x;
      m_values["mainFrame"]["position"]["y"]  = wxDefaultPosition.y;
      m_values["mainFrame"]["size"]["width"]  = wxDefaultSize.GetWidth();
      m_values["mainFrame"]["size"]["height"] = wxDefaultSize.GetHeight();
   }

   void jSDRConfig::SetDefaultConfigValues() { SetDefaultDisplayValues(); }

   auto jSDRConfig::GetDisplayProperties() -> std::shared_ptr<jSDRConfig::DisplayProperties> {
      auto      displayProps          = std::make_shared<jSDRConfig::DisplayProperties>();
      const int x                     = m_values["mainFrame"]["position"]["x"].asInt();   // NOLINT
      const int y                     = m_values["mainFrame"]["position"]["y"].asInt();   // NOLINT
      displayProps->mainFramePosition = wxPoint(x, y);
      const int width                 = m_values["mainFrame"]["size"]["width"].asInt();
      const int height                = m_values["mainFrame"]["size"]["width"].asInt();
      displayProps->mainFrameSize     = wxSize(width, height);
      return displayProps;
   }
}   // namespace jsdr