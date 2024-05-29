export const download = (url: string) => {
  const newTab = window.open();
  if (newTab) {
    newTab.opener = null;
    newTab.location = url;
  } else {
    console.error("Failed to open new tab");
  }
};
