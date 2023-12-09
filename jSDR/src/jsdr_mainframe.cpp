#include "jsdr_mainframe.h"

namespace jsdr {
   jSDRMainFrame::jSDRMainFrame(const wxPoint upperLeft, const wxSize size)
     : wxFrame(nullptr, wxID_ANY, "jSDR", upperLeft, size) {}
}   // namespace jsdr