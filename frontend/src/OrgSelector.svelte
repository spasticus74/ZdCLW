<script>
    import { onMount } from "svelte";
    import { GetOrgNames } from "../wailsjs/go/main/App";

    export let siteOptions = [];
    export let display_func = a => a;
    export let index = 0;
    export let value;

    onMount(async () => {
        let orgs = await GetOrgNames()
        let so = ['Select an organisation...'];
        orgs.forEach(org => {
            so.push(org)
        });
        siteOptions = so
    })
  
    $: {
      value = siteOptions[index];
      console.log(value);
    }
  </script>
  
  <select bind:value={index}>
    {#each siteOptions as siteOption, i}
      <option value={i}>{display_func(siteOption)}</option>
    {/each}
  </select>
  