<script>
  import Tabs from "./Tabs.svelte";
  import contactsTab from "./ContactsTab.svelte";
  import AddCallTab from "./AddCallTab.svelte";

  // List of tab items with labels, values and assigned components
  let items = [
    { label: "Add Call",
      value: 1,
      component: AddCallTab
    },
    { label: "Contacts",
      value: 2,
      component: contactsTab
    }
  ];
</script>

<Tabs {items} />