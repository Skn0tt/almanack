import { reactive, computed, ref, toRefs } from "@vue/composition-api";
import ArcArticle from "./arc-article.js";

export function makeFeed(service) {
  const feed = ref(null);

  const apiState = reactive({
    isLoading: false,
    error: null,
    articles: computed(() =>
      apiState.isLoading || apiState.error || !feed.value
        ? { planned: [], available: [] }
        : ArcArticle.from(feed.value)
    ),
    didLoad: false,
    canLoad: service.hasAuthAvailable(),
  });

  let methods = {
    articleFromID(idFn) {
      return computed(() => {
        let id = idFn();
        return apiState.articles.find(article => article.id === id);
      });
    },
    async reload({ force = false } = {}) {
      if (apiState.isLoading && !force) {
        return;
      }
      apiState.isLoading = true;
      [feed.value, apiState.error] = await service.listAvailable();
      apiState.isLoading = false;
      apiState.didLoad = true;
    },
    async initLoad() {
      if (apiState.canLoad && !apiState.didLoad) {
        await methods.reload();
      }
    },
  };

  return {
    ...toRefs(apiState),
    ...methods,
  };
}
