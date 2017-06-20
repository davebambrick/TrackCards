FROM kylegrantlucas/bazel

RUN git config --global url."https://c91a86f074d53befe091e33fbdc8c17871256f63:x-oauth-basic@github.com/".insteadOf "https://github.com/"
ENV PORT=8080
ENV GO_ENV=development
ADD . /go/src/github.com/iAmPlus/TrackCards
VOLUME /go/.cache/bazel
RUN  cd /go/src/github.com/iAmPlus/TrackCards && bazel build :TrackCards
ENTRYPOINT /go/src/github.com/iAmPlus/TrackCards/bazel-bin/TrackCards

EXPOSE 8080
