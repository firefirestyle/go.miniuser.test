import 'package:firefirestyle.cl.netbox/netbox.dart' as nbox;
import 'package:firefirestyle.httprequest/request.dart' as nbox;
import 'package:firefirestyle.httprequest/request_ver_html.dart' as nbox;
import 'package:firefirestyle.location/location_html.dart' as loc;
//
import 'dart:html' as html;

const String CONFIG_BACKEND_ADDR = "http://localhost:8080";

String accessToken = "";
String userName = "";

loc.HtmlLocation currentLocation = new loc.HtmlLocation();
nbox.UserNBox userNbox = new nbox.UserNBox(new nbox.Html5NetBuilder(), CONFIG_BACKEND_ADDR);
nbox.MeNBox meNbox = new nbox.MeNBox(new nbox.Html5NetBuilder(), CONFIG_BACKEND_ADDR);

main() async {
  print("Hello World");
  if (currentLocation.hashPath == "#/SNS") {
    accessToken = currentLocation.getValueAsString("token", "");
    userName = currentLocation.getValueAsString("userName", "");
//    html.window.location.replace(currentLocation.baseAddr);
  }
  if (accessToken == "") {
    await showLogin();
  } else {
    await showUserInfo(userName);
    await showUserEdit(userName);
  }
}

//
//
//

showLogin() async {
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

showUserInfo(String userName) async {
  nbox.UserInfoProp userInfo = await userNbox.getUserInfo(userName);
  html.document.body.children.add(new html.Element.html(
      [
        """<div>""",
        """<div>DisplayName : ${userInfo.displayName}</div>""", //
        """<br>""",
        """<div>UserName : ${userInfo.userName}</div>""", //
        """<br>""",
        """<div>Created : ${userInfo.created}</div>""", //
        """<div>Point : ${userInfo.point}</div>""", //
        """<div>Content : ${userInfo.content}</div>""", //
        """<div>Sign : ${userInfo.sign}</div>""", //
        """<div>PublicInfo : ${userInfo.publicInfo}</div>""", //
        """<div>PrivateInfo : ${userInfo.privateInfo}</div>""", //
        """<div>IconUrl : ${userInfo.iconUrl}</div>""", //
        """<br>""",
        """</div>""",
      ].join("\r\n"), //
      treeSanitizer: html.NodeTreeSanitizer.trusted));
}


showUserEdit(String userName) async {
  nbox.UserInfoProp userInfo = await userNbox.getUserInfo(userName);
  var cont = new html.Element.html("""<div></div>""");
  html.InputElement displaynameElm = new html.Element.html("""<input type="text" value="${userInfo.displayName}" placeholder="display name">""",treeSanitizer: html.NodeTreeSanitizer.trusted);
  html.InputElement contentElm = new html.Element.html("""<input type="text" value="${userInfo.content}" placeholder="content">""",treeSanitizer: html.NodeTreeSanitizer.trusted);
  var buttonElm = new html.Element.html("""<button>Update</button>""");
  cont.children.add(displaynameElm);
  cont.children.add(contentElm);
  cont.children.add(buttonElm);
  html.document.body.children.add(cont);

  buttonElm.onClick.listen((ev) async {
   nbox.UserInfoProp nextUser = await meNbox.updateUserInfo(accessToken, userName,displayName: displaynameElm.value, cont: contentElm.value);
  });
}
