#include "jSDRApp.h"

#include "jSDRMainFrame.h"

wxIMPLEMENT_APP( jSDRApp );

bool jSDRApp::OnInit() {
   auto           displayProperties = m_config.getDisplayProperties();
   jSDRMainFrame* frame =
       new jSDRMainFrame( displayProperties->mainFramePosition, displayProperties->mainFrameSize );
   frame->Show();
   return true;
}
