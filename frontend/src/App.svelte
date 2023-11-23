<script>
  import Tabs from "./Tabs.svelte";
  //import Tab1 from "./Tab1.svelte";
  import contactsTab from "./ContactsTab.svelte";

  // List of tab items with labels, values and assigned components
  let items = [
    //{ label: "Calls",
    //  value: 1,
    //  component: Tab1
    //},
    { label: "Contacts",
      value: 1,
      component: contactsTab
    }
  ];
</script>

<Tabs {items} />