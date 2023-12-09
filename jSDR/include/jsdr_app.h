#include <wx/wx.h>

#include "jsdr_config.h"

namespace jsdr {
   class jSDRApp : public wxApp {
   public:
      bool        OnInit() override;
      jSDRConfig& Config() { return m_config; }

   private:
      jSDRConfig m_config;
   };
}   // namespace jsdr