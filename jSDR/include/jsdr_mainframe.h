#ifndef JSDR_MAINFRAME_H
#define JSDR_MAINFRAME_H

#include <wx/defs.h>
#include <wx/frame.h>

namespace jsdr {
   class JSdrMainFrame : public wxFrame {
   public:
      JSdrMainFrame(wxPoint upperLeft, wxSize size);
   };
}   // namespace jsdr

#endif   // JSDR_MAINFRAME_H