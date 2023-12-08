#include "jSDRMainFrame.h"

namespace jsdr {
   jSDRMainFrame::jSDRMainFrame(const wxPoint upperLeft, const wxSize size)
     : wxFrame(nullptr, wxID_ANY, "jSDR", upperLeft, size) {}
}   // namespace jsdr