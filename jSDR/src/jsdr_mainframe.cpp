#include "jsdr_mainframe.h"

#include <wx/defs.h>
#include <wx/frame.h>
#include <wx/gdicmn.h>

namespace jsdr {
   JSdrMainFrame::JSdrMainFrame(wxPoint upperLeft, wxSize size)
      : wxFrame(nullptr, wxID_ANY, "jSDR", upperLeft, size) {}
}   // namespace jsdr