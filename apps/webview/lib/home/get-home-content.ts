import homeData from "@/data/webview-home.json";
import type { HomeContent, NavigationItem } from "@/types/home";

type RawNavigationItem = (typeof homeData.header.navigation)[number];

function normalizeNavigationItem(item: RawNavigationItem): NavigationItem {
  if ("children" in item) {
    return { kind: "dropdown", label: item.label, children: item.children ?? [] };
  }

  return { kind: "link", label: item.label, href: item.href };
}

const homeContent: HomeContent = {
  ...homeData,
  header: {
    ...homeData.header,
    navigation: homeData.header.navigation.map(normalizeNavigationItem),
  },
};

export function getHomeContent(): HomeContent {
  return homeContent;
}
