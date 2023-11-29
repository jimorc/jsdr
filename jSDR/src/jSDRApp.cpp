#include "jSDRApp.h"
#include "jSDRMainFrame.h"

wxIMPLEMENT_APP(jSDRApp);

bool jSDRApp::OnInit() {
    jSDRMainFrame *frame = new jSDRMainFrame();
    frame->Show();
    return true;
}

