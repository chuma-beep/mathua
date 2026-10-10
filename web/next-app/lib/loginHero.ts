// The login page's hero picture.
//
// It lives here rather than inline in the route because a Next.js page may only
// export a fixed set of names, and the dimensions need to be readable by a test
// that checks them against the file on disk. FOOTER_ART sets the precedent.

export interface LoginHero {
  /** Public path, as written into the element's src. */
  file: string
  /** The file's real pixel dimensions. */
  width: number
  height: number
}

// A 736x1051 portrait that fills the full height of its column, whatever height
// that column turns out to be. The column is a grid item that stretches to the
// row, and the row is the taller of the viewport and the form beside it — so on a
// short window the picture is taller than the screen and scrolls with the form,
// which is the intended reading of "full height", not a defect.
//
// object-cover, not object-fill: filling would stretch the artwork to whatever
// the column happens to be, and a drawing stretched by even ten percent looks
// wrong in a way that cropping does not. The trade is that the picture is 0.70
// wide-to-tall, so a column of any other ratio loses a slice — off the sides when
// the column is relatively wider, off the top and bottom when it is relatively
// taller. At the widths and heights this is checked at, that is a few percent,
// and the subject sits in the middle of the frame where the slice does not reach.
//
// test/loginHero.test.tsx reads the file and asserts these numbers match, because
// the symptom of a wrong width/height — art that reserves the wrong box before
// the bytes land — only shows up on a slow connection.
export const LOGIN_HERO: LoginHero = {
  file: '/hero/manreader.jpeg',
  width: 736,
  height: 1051,
}