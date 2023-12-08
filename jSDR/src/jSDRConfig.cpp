#include <jSDRConfig.h>

#include <filesystem>
#include <fstream>

jSDRConfig::jSDRConfig() {
   std::filesystem::path cfgFile(getenv("HOME"));
   cfgFile /= m_configFileName;
   if (std::filesystem::exists(cfgFile)) {
      std::ifstream configStream;
      configStream.imbue(std::locale());
      configStream.open(cfgFile);
      configStream >> m_values;
      configStream.close();
   } else {
      setDefaultConfigValues();
   }
}

jSDRConfig::~jSDRConfig() noexcept {
   std::filesystem::path cfgFile(getenv("HOME"));
   cfgFile /= m_configFileName;
   std::ofstream configStream;
   configStream.imbue(std::locale());
   configStream.open(cfgFile);
   configStream << m_values << std::endl;
   configStream.close();
}

void jSDRConfig::setDefaultDisplayValues() {
   m_values["mainFrame"]["displayNumber"]  = 0;
   m_values["mainFrame"]["position"]["x"]  = wxDefaultPosition.x;
   m_values["mainFrame"]["position"]["y"]  = wxDefaultPosition.y;
   m_values["mainFrame"]["size"]["width"]  = wxDefaultSize.GetWidth();
   m_values["mainFrame"]["size"]["height"] = wxDefaultSize.GetHeight();
}

void jSDRConfig::setDefaultConfigValues() { setDefaultDisplayValues(); }

std::shared_ptr<jSDRConfig::DisplayProperties> jSDRConfig::getDisplayProperties() {
   auto displayProps               = std::make_shared<jSDRConfig::DisplayProperties>();
   displayProps->displayNumber     = m_values["mainFrame"]["displayNumber"].asInt();
   int x                           = m_values["mainFrame"]["position"]["x"].asInt();
   int y                           = m_values["mainFrame"]["position"]["y"].asInt();
   displayProps->mainFramePosition = wxPoint(x, y);
   int width                       = m_values["mainFrame"]["size"]["width"].asInt();
   int height                      = m_values["mainFrame"]["size"]["width"].asInt();
   displayProps->mainFrameSize     = wxSize(width, height);
   return displayProps;
}