#include "jsdr_app.h"

#include <wx/app.h>

#include "jsdr_mainframe.h"

wxIMPLEMENT_APP(jsdr::jSDRApp);

namespace jsdr {

   auto jSDRApp::OnInit() -> bool {
      auto           displayProperties = m_config.GetDisplayProperties();
      jSDRMainFrame* frame             =   // NOLINT
          new jSDRMainFrame(displayProperties->mainFramePosition,
                            displayProperties->mainFrameSize);   // NOLINT
      frame->Show();
      return true;
   }
}   // namespace jsdr
