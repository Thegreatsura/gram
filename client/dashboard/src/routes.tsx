import { Icon, IconName, IconProps } from "@speakeasy-api/moonshine";
import { useMemo } from "react";
import { Link, useLocation, useNavigate } from "react-router";
import { useSlugs } from "./contexts/Sdk";
import EnvironmentPage from "./pages/environments/Environment";
import Environments, {
  EnvironmentsRoot,
} from "./pages/environments/Environments";
import Integrations from "./pages/integrations/Integrations";
import Login from "./pages/login/Login";
import { MCPDetailPage } from "./pages/mcp/MCPDetails";
import { MCPOverview, MCPRoot } from "./pages/mcp/MCPOverview";
import Onboarding from "./pages/onboarding/Onboarding";
import { OnboardingWizard } from "./pages/onboarding/Wizard";
import OpenAPIDocuments from "./pages/openapi/OpenAPI";
import Playground from "./pages/playground/Playground";
import NewPromptPage from "./pages/prompts/NewPrompt";
import PromptPage from "./pages/prompts/Prompt";
import Prompts, { PromptsRoot } from "./pages/prompts/Prompts";
import SDK from "./pages/sdk/SDK";
import Settings from "./pages/settings/Settings";
import SlackApp from "./pages/slackapp/SlackApp";
import CustomTools, { CustomToolsRoot } from "./pages/toolBuilder/CustomTools";
import {
  ToolBuilderNew,
  ToolBuilderPage,
} from "./pages/toolBuilder/ToolBuilder";
import { ToolSelect } from "./pages/toolsets/ToolSelect";
import ToolsetPage, { ToolsetRoot } from "./pages/toolsets/Toolset";
import Toolsets, { ToolsetsRoot } from "./pages/toolsets/Toolsets";

type AppRouteBasic = {
  title: string;
  url: string;
  external?: boolean;
  icon?: IconName;
  component?: React.ComponentType;
  indexComponent?: React.ComponentType;
  subPages?: AppRoutesBasic;
  unauthenticated?: boolean;
};

type GoToFunction = (...params: string[]) => void;

export type AppRoutes = Record<string, AppRoute>;
type AppRoutesBasic = Record<string, AppRouteBasic>;

// App route augmented with some additional utilities
export type AppRoute = Omit<AppRouteBasic, "icon" | "subPages"> & {
  Icon: React.ComponentType<Omit<IconProps, "name">>;
  active: boolean;
  // subPages?: AppRoutes;
  href: (...params: string[]) => string;
  goTo: GoToFunction;
  Link: React.ComponentType<{
    params?: string[];
    queryParams?: Record<string, string>;
    children: React.ReactNode;
  }>;
};

type RouteEntry = {
  title: string;
  url: string;
  icon?: IconName;
} & (
  | {
      external: true;

      component?: never;
      indexComponent?: never;
      unauthenticated?: never;
      subPages?: never;
    }
  | {
      external?: false;

      component?: React.ComponentType;
      indexComponent?: React.ComponentType;
      unauthenticated?: boolean;
      subPages?: Record<string, RouteEntry>;
    }
);

const ROUTE_STRUCTURE = {
  login: {
    title: "Login",
    url: "/login",
    component: Login,
    unauthenticated: true,
  },
  onboarding: {
    title: "Onboarding",
    url: ":orgSlug/:projectSlug/onboarding", // Route is like this to break us out of the normal page structure
    component: OnboardingWizard,
  },
  openapi: {
    title: "Your APIs",
    url: "",
    icon: "file-json-2",
    component: OpenAPIDocuments,
  },
  playground: {
    title: "Playground",
    url: "playground",
    icon: "message-circle",
    component: Playground,
  },
  integrations: {
    title: "Integrations",
    url: "integrations",
    icon: "package",
    component: Integrations,
  },
  customTools: {
    title: "Custom Tools",
    url: "custom-tools",
    icon: "pencil-ruler",
    component: CustomToolsRoot,
    indexComponent: CustomTools,
    subPages: {
      toolBuilderNew: {
        title: "Tool Builder",
        url: "new",
        component: ToolBuilderNew,
      },
      toolBuilder: {
        title: "Tool Builder",
        url: ":toolName",
        component: ToolBuilderPage,
      },
    },
  },
  prompts: {
    title: "Prompts",
    url: "prompts",
    icon: "newspaper",
    component: PromptsRoot,
    indexComponent: Prompts,
    subPages: {
      newPrompt: {
        title: "New Prompt",
        url: "new",
        component: NewPromptPage,
      },
      prompt: {
        title: "Edit Prompt",
        url: ":promptName",
        component: PromptPage,
      },
    },
  },
  toolsets: {
    title: "Toolsets",
    url: "toolsets",
    icon: "blocks",
    component: ToolsetsRoot,
    indexComponent: Toolsets,
    subPages: {
      toolset: {
        title: "Toolset",
        url: ":toolsetSlug",
        component: ToolsetRoot,
        indexComponent: ToolsetPage,
        subPages: {
          update: {
            title: "Update",
            url: "update",
            component: ToolSelect,
          },
        },
      },
    },
  },
  mcp: {
    title: "MCP",
    url: "mcp",
    icon: "network",
    component: MCPRoot,
    indexComponent: MCPOverview,
    subPages: {
      details: {
        title: "MCP Details",
        url: ":toolsetSlug",
        component: MCPDetailPage,
      },
    },
  },
  environments: {
    title: "Environments",
    url: "environments",
    icon: "globe",
    component: EnvironmentsRoot,
    indexComponent: Environments,
    subPages: {
      environment: {
        title: "Environment",
        url: ":environmentSlug",
        component: EnvironmentPage,
      },
    },
  },
  agents: {
    title: "Agents",
    url: "agents",
    icon: "code",
    component: SDK,
  },
  slackApp: {
    title: "Slack App",
    url: "slack-app",
    icon: "slack",
    component: SlackApp,
  },
  uploadOpenAPI: {
    title: "Upload OpenAPI",
    url: "upload",
    icon: "upload",
    component: Onboarding,
  },
  settings: {
    title: "Settings",
    url: "settings",
    icon: "settings",
    component: Settings,
  },
  docs: {
    title: "Docs",
    url: "https://docs.getgram.ai",
    icon: "book-open",
    external: true,
  },
} satisfies Record<string, RouteEntry>;

type RouteStructure = typeof ROUTE_STRUCTURE;

/**
 * The point of all this type magic is to make it so you only have to define the routes once
 * and the `useRoutes` hook can add a lot of extra utilities without losing the type safety.
 */

// Transform the AppRouteBasic into an AppRoute, recursing on subPages if present
// so that subPages keeps its route-specific type
type TransformAppRoute<T extends AppRouteBasic> = T extends {
  subPages: AppRoutesBasic;
}
  ? Omit<AppRoute, "subPages"> & TransformRouteToGoTo<T["subPages"]>
  : AppRoute;

type TransformElem<T> = T extends AppRouteBasic
  ? TransformAppRoute<T>
  : T extends AppRouteBasic
  ? TransformRouteToGoTo<T>
  : T;

type TransformRouteToGoTo<T> = {
  [K in keyof T]: TransformElem<T[K]>;
};

type RoutesWithGoTo = TransformRouteToGoTo<RouteStructure>;

export const useRoutes = (): RoutesWithGoTo => {
  const location = useLocation();
  const { orgSlug, projectSlug } = useSlugs();
  const navigate = useNavigate();

  // Check if the current url matches the route url, including dynamic segments
  const matchesCurrent = (url: string) => {
    const urlParts = url.split("/").filter(Boolean);
    const currentParts = location.pathname.split("/").filter(Boolean);

    if (urlParts.length !== currentParts.length) {
      return false;
    }

    return urlParts.every(
      (part, index) => part === currentParts[index] || part.startsWith(":")
    );
  };

  const addRouteUtilities = (
    route: AppRouteBasic,
    parent?: string
  ): AppRoute => {
    if (parent === undefined && !route.url.startsWith("/")) {
      parent = `/:orgSlug/:projectSlug`;
    }

    const urlWithParent = `${parent ?? ""}/${route.url}`;

    const resolveUrl = (...params: string[]) => {
      if (route.external) {
        return route.url;
      }

      const parts = urlWithParent.split("/").filter(Boolean);
      const finalParts = [];

      for (const part of parts) {
        if (part.startsWith(":")) {
          if (part === ":orgSlug") {
            finalParts.push(orgSlug);
          } else if (part === ":projectSlug") {
            finalParts.push(projectSlug);
          } else {
            const v = params.shift();
            if (!v) {
              throw new Error(`No value provided for ${part}`);
            }
            finalParts.push(v);
          }
        } else {
          finalParts.push(part);
        }
      }

      return "/" + finalParts.join("/");
    };

    const goTo = (...params: string[]) => {
      navigate(resolveUrl(...params));
    };

    const linkComponent = ({
      params = [],
      queryParams = {},
      children,
    }: {
      params?: string[];
      queryParams?: Record<string, string>;
      children: React.ReactNode;
    }) => {
      const queryString = new URLSearchParams(queryParams).toString();
      return (
        <Link
          to={`${resolveUrl(...params)}?${queryString}`}
          className="hover:underline"
        >
          {children}
        </Link>
      );
    };

    const subPages = route.subPages
      ? addGoToToRoutes(route.subPages, urlWithParent)
      : undefined;

    const active =
      matchesCurrent(urlWithParent) ||
      !!Object.values(subPages ?? {}).some((subPage) => subPage.active);

    const newRoute: AppRoute = {
      ...route,
      active,
      Icon: (props: Omit<IconProps, "name">) =>
        route.icon ? <Icon {...props} name={route.icon} /> : null,
      href: resolveUrl,
      goTo,
      Link: linkComponent,
      ...subPages,
    };

    if (route.url.startsWith("/")) {
      newRoute.goTo = () => route.url;
    }

    return newRoute;
  };

  const addGoToToRoutes = <T extends AppRoutesBasic>(
    routes: T,
    parent?: string
  ): TransformRouteToGoTo<T> => {
    return Object.fromEntries(
      Object.entries(routes).map(([key, route]) => [
        key,
        addRouteUtilities(route, parent),
      ])
    ) as TransformRouteToGoTo<T>;
  };

  const routes: RoutesWithGoTo = useMemo(
    () => addGoToToRoutes(ROUTE_STRUCTURE),
    [location.pathname]
  );

  return routes;
};
