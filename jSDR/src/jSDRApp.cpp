#include "jSDRApp.h"

#include "jSDRMainFrame.h"

wxIMPLEMENT_APP(jsdr::jSDRApp);

namespace jsdr {

   bool jSDRApp::OnInit() {
      auto           displayProperties = m_config.getDisplayProperties();
      jSDRMainFrame* frame             =   // NOLINT
          new jSDRMainFrame(displayProperties->mainFramePosition,
                            displayProperties->mainFrameSize);   // NOLINT
      frame->Show();
      return true;
   }
}   // namespace jsdr
