<script>
  import { navigate } from "svelte-routing";
  import { IsUserAdmin } from "../client/user";
  import { onMount } from "svelte";

  onMount(async () => {
    let isAdmin = await IsUserAdmin();
    if (isAdmin) {
      navigate("/admin");
    } else {
      navigate("/public-index");
    }
  });
</script>

<div></div>
