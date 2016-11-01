import 'package:firefirestyle.cl.netbox/netbox.dart' as nbox;
import 'package:firefirestyle.httprequest/request.dart' as nbox;
import 'package:firefirestyle.httprequest/request_ver_html.dart' as nbox;
import 'package:firefirestyle.location/location_html.dart' as loc;
//
import 'dart:html' as html;

const String CONFIG_BACKEND_ADDR = "http://localhost:8080";

String accessToken = "";
String userName = "";
main() {
  print("Hello World");
  nbox.UserNBox userNbox = new nbox.UserNBox(new nbox.Html5NetBuilder(), CONFIG_BACKEND_ADDR);
  nbox.MeNBox meNbox = new nbox.MeNBox(new nbox.Html5NetBuilder(), CONFIG_BACKEND_ADDR);

  loc.HtmlLocation currentLocation = new loc.HtmlLocation();
  if(currentLocation.hashPath == "#/SNS" ) {
    accessToken = currentLocation.getValueAsString("token", "");
    userName = currentLocation.getValueAsString("userName", "");
    html.window.location.replace(currentLocation.baseAddr);
  }
  html.document.body.children.add(new html.Element.html(
      [
        """<div>""",
        """<a href="${meNbox.makeLoginTwitterUrl(currentLocation.baseAddr+"/#/SNS")}">Twitter Login</a>""", //
        """<br>""",
        """<a href="${meNbox.makeLoginFacebookUrl(currentLocation.baseAddr+"/#/SNS")}">Facebook Login</a>""" //
        """</div>""",
      ].join("\r\n"), //
      treeSanitizer: html.NodeTreeSanitizer.trusted));
}
