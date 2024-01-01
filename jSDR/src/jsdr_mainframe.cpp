#include "jsdr_mainframe.h"

#include <wx/defs.h>
#include <wx/frame.h>
#include <wx/gdicmn.h>

namespace jsdr {
   jSDRMainFrame::jSDRMainFrame(const wxPoint upperLeft, const wxSize size)
      : wxFrame(nullptr, wxID_ANY, "jSDR", upperLeft, size) {}
}   // namespace jsdr