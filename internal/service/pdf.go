package service

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type PDFService struct {
	workerSem chan struct{}
}

func NewPDFService(maxConcurrent uint) *PDFService {
	return &PDFService{
		workerSem: make(chan struct{}, maxConcurrent),
	}
}

func (s *PDFService) GenerateFromHTML(ctx context.Context, html string) ([]byte, error) {
	// 1. Acquire a worker slot (blocks if at maxConcurrent)
	s.workerSem <- struct{}{}
	defer func() { <-s.workerSem }()

	// 2. Setup chromedp context with a timeout
	// Always use a timeout so a "zombie" process doesn't eat your RAM forever
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,  // Required for some Linux/Docker environments
		chromedp.Headless,   // Don't try to open a window
		chromedp.DisableGPU, // Saves memory
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	defer taskCancel()

	var buf []byte
	err := chromedp.Run(taskCtx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			tree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(tree.Frame.ID, html).Do(ctx)
		}),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(800*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				Do(ctx)
			return err
		}),
	)

	return buf, err
}
