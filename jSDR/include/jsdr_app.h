#ifndef JSDR_APP_H
#define JSDR_APP_H

#include <wx/wx.h>

#include "jsdr_config.h"

namespace jsdr {
   class JSdrApp : public wxApp {
   public:
      auto OnInit() -> bool override;
      auto Config() -> JSdrConfig& { return _config; }

   private:
      JSdrConfig _config;
   };
}   // namespace jsdr

#endif   // JSDR_APP_H