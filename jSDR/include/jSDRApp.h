#include <wx/wx.h>
#include <json/json.h>

class jSDRApp : public wxApp {
public:
    bool OnInit() override;
private:
    Json::Value m_root;
};
