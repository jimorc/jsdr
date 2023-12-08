#include <json/json.h>
#include <wx/wx.h>

#include <memory>
#include <string>

namespace jsdr {
   class jSDRConfig {
   public:
      struct DisplayProperties {
         int     displayNumber     = 0;
         wxPoint mainFramePosition = wxDefaultPosition;
         wxSize  mainFrameSize     = wxDefaultSize;
      };
      jSDRConfig();
      ~jSDRConfig() noexcept;
      std::shared_ptr<DisplayProperties> getDisplayProperties();
      Json::Value&                       values() { return m_values; }

   private:
      void        setDefaultConfigValues();
      void        setDefaultDisplayValues();
      std::string m_configFileName{ "jSDR.config" };
      Json::Value m_values;
   };
}   // namespace jsdr