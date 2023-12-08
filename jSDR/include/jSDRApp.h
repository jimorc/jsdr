#include <wx/wx.h>

#include "jSDRConfig.h"

class jSDRApp : public wxApp {
public:
   bool        OnInit() override;
   jSDRConfig& Config() { return m_config; }

private:
   jSDRConfig m_config;
};
