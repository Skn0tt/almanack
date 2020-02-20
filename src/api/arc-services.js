import { computed } from "@vue/composition-api";
import ArcArticle from "./arc-article.js";
import { useService } from "./service.js";

export function listAvailable(service) {
  let apiState = useService({
    canLoad: service.hasAuthAvailable(),
    serviceCall: () => service.listAvailable(),
  });

  return {
    ...apiState,
    articles: computed(() =>
      apiState.isLoading.value ||
      apiState.error.value ||
      !apiState.rawData.value
        ? { planned: [], available: [] }
        : ArcArticle.from(apiState.rawData.value)
    ),
  };
}
