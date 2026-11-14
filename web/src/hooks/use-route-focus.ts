import React from "react"
import { useLocation, useNavigationType } from "react-router-dom"

export const PAGE_TITLE_ID = "page-title"
export const MAIN_CONTENT_ID = "main-content"

const OPEN_DIALOG =
  '[role="dialog"][data-state="open"], [role="alertdialog"][data-state="open"]'

const pageHasFocus = () => {
  const main = document.getElementById(MAIN_CONTENT_ID)
  const active = document.activeElement
  return !!main && active !== main && main.contains(active)
}

// useRouteFocus moves focus to the page heading on a move to another page, instead of leaving it on the link that left.
export function useRouteFocus() {
  const { pathname } = useLocation()
  const navigationType = useNavigationType()
  const previous = React.useRef(pathname)

  React.useEffect(() => {
    if (previous.current === pathname) return
    previous.current = pathname
    // A redirect is no move of the user's; an open dialog or a page's own focus wins.
    if (
      navigationType === "REPLACE" ||
      document.querySelector(OPEN_DIALOG) ||
      pageHasFocus()
    ) {
      return
    }
    document.getElementById(PAGE_TITLE_ID)?.focus({ preventScroll: true })
  }, [pathname, navigationType])
}
