#include "jsdr_app.h"

#include <wx/app.h>

#include "jsdr_mainframe.h"

wxIMPLEMENT_APP(jsdr::JSdrApp);

namespace jsdr {

   auto JSdrApp::OnInit() -> bool {
      auto           displayProperties = _config.GetDisplayProperties();
      JSdrMainFrame* frame             =   // NOLINT
          new JSdrMainFrame(displayProperties->mainFramePosition, displayProperties->mainFrameSize);
      frame->Show();
      return true;
   }
}   // namespace jsdr
